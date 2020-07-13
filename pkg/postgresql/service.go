package postgresql

import (
	postgresqlv1alpha1 "postgres-operator/pkg/apis/postgresql/v1alpha1"
	"postgres-operator/pkg/utils"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//NewServiceHeadlessForCR creates a new Headless Service
func NewServiceHeadlessForCR(cr *postgresqlv1alpha1.PostgreSQL) *v1.Service {
	labels := utils.LabelsForPostgreSQL(cr.ObjectMeta.Name)
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: cr.ObjectMeta.Name + "-postgresql-headless", Labels: labels, Namespace: cr.ObjectMeta.Namespace},
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

//NewServiceForCR creates a new Service
func NewServiceForCR(cr *postgresqlv1alpha1.PostgreSQL) *v1.Service {
	labels := utils.LabelsForPostgreSQL(cr.ObjectMeta.Name)
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: cr.ObjectMeta.Name + "-postgresql", Labels: labels, Namespace: cr.ObjectMeta.Namespace},
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
