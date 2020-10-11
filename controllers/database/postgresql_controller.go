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
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/go-logr/logr"
	databasev1alpha1 "github.com/iboware/postgresql-operator/apis/database/v1alpha1"
	postgres "github.com/iboware/postgresql-operator/postgres"
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
	crd := &databasev1alpha1.PostgreSQL{}
	err := r.Get(context.TODO(), req.NamespacedName, crd)
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
	namespace := postgres.NewNamespaceForCR(crd)
	postgresStatefulSet := postgres.NewStatefulSet(crd)
	postgresService := postgres.NewPostgresService(crd)
	postgresServiceHeadless := postgres.NewPostgressHeadlessService(crd)
	postgresSecret := postgres.NewPostgresSecret(crd)
	posgresConfigMap, err := postgres.NewConfigMap(crd)
	pgPoolDeployment := postgres.NewPgPoolDeployment(crd)
	pgPoolService := postgres.NewPgPoolService(crd)

	namespaceFound := &v1.Namespace{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: namespace.Name}, namespaceFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(crd, namespace, r.Scheme)
			err = r.Create(context.TODO(), namespace)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			r.Log.Info("failed to get Namespace")
			return reconcile.Result{}, err
		}
	}

	secretFound := &v1.Secret{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: postgresSecret.Name, Namespace: postgresSecret.Namespace}, secretFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(crd, postgresSecret, r.Scheme)
			err = r.Create(context.TODO(), postgresSecret)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			r.Log.Info("failed to get Postgres Secret")
			return reconcile.Result{}, err
		}
	}

	configMapFound := &v1.ConfigMap{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: posgresConfigMap.Name, Namespace: posgresConfigMap.Namespace}, configMapFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(crd, posgresConfigMap, r.Scheme)
			err = r.Create(context.TODO(), posgresConfigMap)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			r.Log.Info("failed to get Postgres ConfigMap")
			return reconcile.Result{}, err
		}
	}

	serviceFound := &v1.Service{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: postgresService.Name, Namespace: postgresService.Namespace}, serviceFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(crd, postgresService, r.Scheme)
			err = r.Create(context.TODO(), postgresService)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			r.Log.Info("failed to get Postgres Service")
			return reconcile.Result{}, err
		}
	}

	serviceHeadlessFound := &v1.Service{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: postgresServiceHeadless.Name, Namespace: postgresServiceHeadless.Namespace}, serviceHeadlessFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(crd, postgresServiceHeadless, r.Scheme)
			err = r.Create(context.TODO(), postgresServiceHeadless)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			r.Log.Info("failed to get Postgres Headless Service")
			return reconcile.Result{}, err
		}
	}

	statefulSetFound := &appsv1.StatefulSet{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: postgresStatefulSet.Name, Namespace: postgresStatefulSet.Namespace}, statefulSetFound)
	if err != nil {
		if errors.IsNotFound(err) {
			controllerutil.SetControllerReference(crd, postgresStatefulSet, r.Scheme)
			err = r.Create(context.TODO(), postgresStatefulSet)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			r.Log.Info("failed to get Postgres StatefulSet")
			return reconcile.Result{}, err
		}
	} else if *statefulSetFound.Spec.Replicas != crd.Spec.Replicas {
		var originalStatefulSet, err1 = json.Marshal(statefulSetFound)
		if err1 != nil {
			return reconcile.Result{}, err1
		}
		statefulSetFound.Spec.Replicas = &crd.Spec.Replicas
		statefulSetFound.Spec.Template.Spec.Containers[0].Env = postgres.EnvVarsForStatefulSet(crd.ObjectMeta.Name, int(crd.Spec.Replicas), crd.Spec.Namespace, crd.Spec.EnablePgPool)

		var newStatefulSet, err2 = json.Marshal(statefulSetFound)
		if err2 != nil {
			return reconcile.Result{}, err2
		}

		var patchBytes, err3 = jsonpatch.CreateMergePatch(originalStatefulSet, newStatefulSet)
		if err != nil {
			return reconcile.Result{}, err3
		}

		err = r.Patch(context.TODO(), statefulSetFound, client.RawPatch(types.MergePatchType, patchBytes))
		if err != nil {
			return reconcile.Result{}, err
		}
		r.Log.Info("Postgres StatefulSet updated")
	}

	//PgPool
	if crd.Spec.EnablePgPool {
		pgPoolServiceFound := &v1.Service{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: pgPoolService.Name, Namespace: pgPoolService.Namespace}, pgPoolServiceFound)
		if err != nil {
			if errors.IsNotFound(err) {
				controllerutil.SetControllerReference(crd, pgPoolService, r.Scheme)
				err = r.Create(context.TODO(), pgPoolService)
				if err != nil {
					return reconcile.Result{}, err
				}
			} else {
				r.Log.Info("failed to get pgPool Service")
				return reconcile.Result{}, err
			}
		}

		pgPoolDeploymentFound := &appsv1.Deployment{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: pgPoolDeployment.Name, Namespace: pgPoolDeployment.Namespace}, pgPoolDeploymentFound)
		if err != nil {
			if errors.IsNotFound(err) {
				controllerutil.SetControllerReference(crd, pgPoolDeployment, r.Scheme)
				err = r.Create(context.TODO(), pgPoolDeployment)
				if err != nil {
					return reconcile.Result{}, err
				}
			} else {
				r.Log.Info("failed to get pgPool Deployment")
				return reconcile.Result{}, err
			}
		} else if *pgPoolDeploymentFound.Spec.Replicas != crd.Spec.Replicas {
			var originalDeployment, err1 = json.Marshal(pgPoolDeploymentFound)
			if err1 != nil {
				return reconcile.Result{}, err1
			}
			pgPoolDeploymentFound.Spec.Replicas = &crd.Spec.Replicas
			pgPoolDeploymentFound.Spec.Template.Spec.Containers[0].Env = postgres.EnvVarsForPgPool(crd.ObjectMeta.Name, int(crd.Spec.Replicas), crd.Spec.Namespace)

			var newDeployment, err2 = json.Marshal(pgPoolDeploymentFound)
			if err2 != nil {
				return reconcile.Result{}, err2
			}

			var patchBytes, err3 = jsonpatch.CreateMergePatch(originalDeployment, newDeployment)
			if err != nil {
				return reconcile.Result{}, err3
			}

			err = r.Patch(context.TODO(), pgPoolDeploymentFound, client.RawPatch(types.MergePatchType, patchBytes))
			if err != nil {
				return reconcile.Result{}, err
			}
			r.Log.Info("pgPool Deployment updated")
		}
	}
	r.Status().Update(context.TODO(), crd)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up reconciler.
func (r *PostgreSQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1alpha1.PostgreSQL{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&v1.ConfigMap{}).
		Owns(&v1.Secret{}).
		Owns(&v1.Service{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}
