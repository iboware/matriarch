package postgresql

import (
	"context"
	"reflect"

	postgresqlv1alpha1 "postgres-operator/pkg/apis/postgresql/v1alpha1"
	"postgres-operator/pkg/postgresql"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
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

	// Define new objects
	statefulSet := postgresql.NewStatefulSetForCR(instance)
	service := postgresql.NewServiceForCR(instance)
	serviceHeadless := postgresql.NewServiceHeadlessForCR(instance)
	secret := postgresql.NewSecretForCR(instance)
	configMap, err := postgresql.NewConfigMapForCR(instance)

	secretFound := &v1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, secretFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(instance, secret, r.scheme)
			err = r.client.Create(context.TODO(), secret)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			reqLogger.Info("failed to get secret")
			return reconcile.Result{}, err
		}
	} else if !reflect.DeepEqual(secret.StringData, secretFound.StringData) {
		secret.ObjectMeta = secretFound.ObjectMeta

		controllerutil.SetControllerReference(instance, secret, r.scheme)
		err = r.client.Update(context.TODO(), secret)
		if err != nil {
			return reconcile.Result{}, err
		}
		reqLogger.Info("secret updated")
	}

	configMapFound := &v1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, configMapFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(instance, configMap, r.scheme)
			err = r.client.Create(context.TODO(), configMap)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			reqLogger.Info("failed to get ConfigMap")
			return reconcile.Result{}, err
		}
	} else if !reflect.DeepEqual(configMap.Data, configMapFound.Data) {
		configMap.ObjectMeta = configMapFound.ObjectMeta
		controllerutil.SetControllerReference(instance, configMap, r.scheme)
		err = r.client.Update(context.TODO(), configMap)
		if err != nil {
			return reconcile.Result{}, err
		}
		reqLogger.Info("ConfigMap updated")
	}

	serviceFound := &v1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, serviceFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(instance, service, r.scheme)
			err = r.client.Create(context.TODO(), service)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			reqLogger.Info("failed to get Service")
			return reconcile.Result{}, err
		}
	} else if !reflect.DeepEqual(service.Spec, serviceFound.Spec) {
		service.ObjectMeta = serviceFound.ObjectMeta
		controllerutil.SetControllerReference(instance, service, r.scheme)
		err = r.client.Update(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}
		reqLogger.Info("Service updated")
	}

	serviceHeadlessFound := &v1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: serviceHeadless.Name, Namespace: serviceHeadless.Namespace}, serviceHeadlessFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(instance, serviceHeadless, r.scheme)
			err = r.client.Create(context.TODO(), serviceHeadless)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			reqLogger.Info("failed to get Headless Service")
			return reconcile.Result{}, err
		}
	} else if !reflect.DeepEqual(serviceHeadless.Spec, serviceHeadlessFound.Spec) {
		serviceHeadless.ObjectMeta = serviceHeadlessFound.ObjectMeta
		controllerutil.SetControllerReference(instance, serviceHeadless, r.scheme)
		err = r.client.Update(context.TODO(), serviceHeadless)
		if err != nil {
			return reconcile.Result{}, err
		}
		reqLogger.Info("Headless Service updated")
	}

	statefulSetFound := &appsv1.StatefulSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, statefulSetFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(instance, statefulSet, r.scheme)
			err = r.client.Create(context.TODO(), statefulSet)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			reqLogger.Info("failed to get statefulSet")
			return reconcile.Result{}, err
		}
	} else if !reflect.DeepEqual(statefulSet.Spec, statefulSetFound.Spec) {
		statefulSet.ObjectMeta = statefulSetFound.ObjectMeta
		controllerutil.SetControllerReference(instance, statefulSet, r.scheme)
		err = r.client.Update(context.TODO(), statefulSet)
		if err != nil {
			return reconcile.Result{}, err
		}
		reqLogger.Info("statefulSet updated")
	}

	r.client.Status().Update(context.TODO(), instance)
	return reconcile.Result{}, nil
}
