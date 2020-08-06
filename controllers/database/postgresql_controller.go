/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	databasev1alpha1 "github.com/iboware/postgresql-operator/apis/database/v1alpha1"
	"github.com/iboware/postgresql-operator/postgresql"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// PostgreSQLReconciler reconciles a PostgreSQL object
type PostgreSQLReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=database.iboware.com,resources=postgresqls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=database.iboware.com,resources=postgresqls/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments;pods;daemonsets;replicasets;statefulsets,verbs=get;update;patch;list;create;delete;watch
// +kubebuilder:rbac:groups=core,resources=secrets;configmaps;services;persistentvolumeclaims;namespaces,verbs=get;update;patch;list;create;delete;watch

// Reconcile is a function
func (r *PostgreSQLReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("postgresql", req.Name)

	// Fetch the PostgreSQL instance
	instance := &databasev1alpha1.PostgreSQL{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
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
	namespace := postgresql.NewNamespaceForCR(instance)
	statefulSet := postgresql.NewStatefulSetForCR(instance)
	service := postgresql.NewServiceForCR(instance)
	serviceHeadless := postgresql.NewServiceHeadlessForCR(instance)
	secret := postgresql.NewSecretForCR(instance)
	configMap, err := postgresql.NewConfigMapForCR(instance)

	namespaceFound := &v1.Namespace{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: namespace.Name}, namespaceFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(instance, namespace, r.Scheme)
			err = r.Client.Create(context.TODO(), namespace)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			r.Log.Info("failed to get namespace")
			return reconcile.Result{}, err
		}
	}

	secretFound := &v1.Secret{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, secretFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(instance, secret, r.Scheme)
			err = r.Client.Create(context.TODO(), secret)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			r.Log.Info("failed to get secret")
			return reconcile.Result{}, err
		}
	}

	configMapFound := &v1.ConfigMap{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, configMapFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(instance, configMap, r.Scheme)
			err = r.Client.Create(context.TODO(), configMap)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			r.Log.Info("failed to get ConfigMap")
			return reconcile.Result{}, err
		}
	}

	serviceFound := &v1.Service{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, serviceFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(instance, service, r.Scheme)
			err = r.Client.Create(context.TODO(), service)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			r.Log.Info("failed to get Service")
			return reconcile.Result{}, err
		}
	}

	serviceHeadlessFound := &v1.Service{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: serviceHeadless.Name, Namespace: serviceHeadless.Namespace}, serviceHeadlessFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(instance, serviceHeadless, r.Scheme)
			err = r.Client.Create(context.TODO(), serviceHeadless)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			r.Log.Info("failed to get Headless Service")
			return reconcile.Result{}, err
		}
	}

	statefulSetFound := &appsv1.StatefulSet{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, statefulSetFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(instance, statefulSet, r.Scheme)
			err = r.Client.Create(context.TODO(), statefulSet)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			r.Log.Info("failed to get statefulSet")
			return reconcile.Result{}, err
		}
	} else if (statefulSetFound.Spec.Replicas != &instance.Spec.Replicas) && (statefulSetFound.Status.ReadyReplicas == *statefulSetFound.Spec.Replicas) {
		statefulSetFound.Spec.Replicas = &instance.Spec.Replicas
		err = r.Client.Update(context.TODO(), statefulSetFound)
		if err != nil {
			return reconcile.Result{}, err
		}
		r.Log.Info("statefulSet updated")
	}

	r.Client.Status().Update(context.TODO(), instance)

	return ctrl.Result{}, nil
}

// SetupWithManager is a function
func (r *PostgreSQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1alpha1.PostgreSQL{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&v1.ConfigMap{}).
		Owns(&v1.Secret{}).
		Owns(&v1.Service{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 2}).
		Complete(r)
}
