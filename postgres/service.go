package postgres

import (
	databasev1alpha1 "github.com/iboware/postgresql-operator/apis/database/v1alpha1"
	utils "github.com/iboware/postgresql-operator/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//NewPostgressHeadlessService creates a new headless service
func NewPostgressHeadlessService(cr *databasev1alpha1.PostgreSQL) *v1.Service {
	labels := utils.NewLabels("postgres", cr.ObjectMeta.Name, "postgres")
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.ObjectMeta.Name + "-postgresql-headless",
			Labels:    labels,
			Namespace: cr.Spec.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceTypeClusterIP,
			Selector: labels,
			Ports: []v1.ServicePort{
				{
					Name:       "postgresql",
					TargetPort: intstr.IntOrString{Type: intstr.String, StrVal: "postgresql"},
					Port:       5432,
					Protocol:   v1.ProtocolTCP,
				},
			}},
	}

	return service
}

//NewPostgresService creates a new service
func NewPostgresService(cr *databasev1alpha1.PostgreSQL) *v1.Service {
	labels := utils.NewLabels("postgres", cr.ObjectMeta.Name, "postgres")
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.ObjectMeta.Name + "-postgresql",
			Labels:    labels,
			Namespace: cr.Spec.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceTypeClusterIP,
			Selector: labels,
			Ports: []v1.ServicePort{
				{
					Name:       "postgresql",
					TargetPort: intstr.IntOrString{Type: intstr.String, StrVal: "postgresql"},
					Port:       5432,
					Protocol:   v1.ProtocolTCP,
				},
			}},
	}

	return service
}

//NewPgPoolService creates a new PgPool service
func NewPgPoolService(cr *databasev1alpha1.PostgreSQL) *v1.Service {
	labels := utils.NewLabels("pgpool", cr.ObjectMeta.Name, "pgpool")
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.ObjectMeta.Name + "-pgpool",
			Labels:    labels,
			Namespace: cr.Spec.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceTypeClusterIP,
			Selector: labels,
			Ports: []v1.ServicePort{
				{
					Name:       "postgresql",
					TargetPort: intstr.IntOrString{Type: intstr.String, StrVal: "postgresql"},
					Port:       5432,
					Protocol:   v1.ProtocolTCP,
				},
			}},
	}

	return service
}
