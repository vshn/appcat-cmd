package defaults

import (
	k8up "github.com/k8up-io/k8up/v2/api/v1"
	"github.com/vshn/appcat-cli/internal/util"
	v1 "k8s.io/api/core/v1"
)

func (d *Defaults) GetBackupDefault(input []util.Input) *k8up.Backup {
	var k8upBackupDefault k8up.Backup
	var jobHistoryVal int = 2
	k8upBackupDefault.Spec.FailedJobsHistoryLimit = &jobHistoryVal
	k8upBackupDefault.Spec.SuccessfulJobsHistoryLimit = &jobHistoryVal
	//The 'empty' initializations need to happen -> panic otherwise
	k8upBackupDefault.Spec.Backend = &k8up.Backend{}
	k8upBackupDefault.Spec.Backend.RepoPasswordSecretRef = &v1.SecretKeySelector{}

	util.DecorateType(&k8upBackupDefault, input)
	return &k8upBackupDefault
}

func (d *Defaults) GetScheduleDefault(input []util.Input) *k8up.Schedule {
	var k8upScheduleDefault k8up.Schedule

	util.DecorateType(&k8upScheduleDefault, input)
	return &k8upScheduleDefault
}
