package postgresql

import (
	"context"

	postgresqlv1alpha1 "postgres-operator/pkg/apis/postgresql/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_postgresql")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new PostgreSQL Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePostgreSQL{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("postgresql-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource PostgreSQL
	err = c.Watch(&source.Kind{Type: &postgresqlv1alpha1.PostgreSQL{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner PostgreSQL
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &postgresqlv1alpha1.PostgreSQL{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcilePostgreSQL implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilePostgreSQL{}

// ReconcilePostgreSQL reconciles a PostgreSQL object
type ReconcilePostgreSQL struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a PostgreSQL object and makes changes based on the state read
// and what is in the PostgreSQL.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePostgreSQL) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling PostgreSQL")

	// Fetch the PostgreSQL instance
	instance := &postgresqlv1alpha1.PostgreSQL{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	statefulSet := newStatefulSetForCR(instance)

	// Set PostgreSQL instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, statefulSet, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &appsv1.StatefulSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "statefulSet.Namespace", statefulSet.Namespace, "statefulSet.Name", statefulSet.Name)
		err = r.client.Create(context.TODO(), statefulSet)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: StatefulSet already exists", "StatefulSet.Namespace", found.Namespace, "StatefulSet.Name", found.Name)
	return reconcile.Result{}, nil
}

//[spec.selector: Required value, spec.template.metadata.labels: Invalid value: map[string]string(nil): `selector` does not match template `labels`]"
// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newStatefulSetForCR(cr *postgresqlv1alpha1.PostgreSQL) *appsv1.StatefulSet {
	labels := labelsForPostgreSQL(cr.Name)

	// Constants for hello-stateful StatefulSet & Volumes
	const (
		AppImage         = "nginxdemos/nginx-hello:latest"
		AppContainerName = "hello-stateful"
		ImagePullPolicy  = v1.PullIfNotPresent
		DiskSize         = 1 * 1000 * 1000 * 1000
	)

	var (
		// storageClassName              = "standard"
		// diskSize                      = *resource.NewQuantity(DiskSize, resource.DecimalSI)
		terminationGracePeriodSeconds = int64(10)
		// accessMode                    = []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce}
		// resourceList                  = v1.ResourceList{v1.ResourceStorage: diskSize}
	)

	statefulset := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.ObjectMeta.Name,
			Namespace: cr.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			ServiceName: cr.ObjectMeta.Name,
			Replicas:    &cr.Spec.Size,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					Containers: []v1.Container{
						{
							Name:            AppContainerName,
							Image:           cr.Spec.Image,
							ImagePullPolicy: ImagePullPolicy,
							// VolumeMounts: []v1.VolumeMount{
							// 	{
							// 		Name:      AppVolumeName,
							// 		MountPath: AppVolumeMountPath,
							// 	},
							// },
						},
					},
					// Volumes: []v1.Volume{
					// 	{
					// 		Name: AppVolumeName,
					// 		VolumeSource: v1.VolumeSource{
					// 			PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
					// 				ClaimName: cr.ObjectMeta.Name,
					// 			},
					// 		},
					// 	},
					// },
				},
			},
		},
	}
	return statefulset
}

// labelsForMemcached returns the labels for selecting the resources
// belonging to the given PostgreSQL CR name.
func labelsForPostgreSQL(name string) map[string]string {
	return map[string]string{"app": "postgresql", "postgresql_cr": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []v1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// func newServiceForCR(cr *postgresqlv1alpha1.PostgreSQL) *corev1.Service {

// }
