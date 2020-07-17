package state

import (
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/openshift/osde2e/pkg/common/alert"
	"github.com/openshift/osde2e/pkg/common/helper"
	"github.com/openshift/osde2e/pkg/common/runner"
)

const (
	// cmd to collect prometheus data
	promCollectCmd = "oc exec -n openshift-monitoring prometheus-k8s-0 -c prometheus -- /bin/sh -c \"cp -ruf /prometheus /tmp/prometheus && tar cvzO -C /tmp/prometheus . "
)

var _ = ginkgo.Describe("[Suite: e2e] Cluster state", func() {
	defer ginkgo.GinkgoRecover()
	ginkgo.BeforeEach(func() {
		alert.RegisterGinkgoAlert(ginkgo.CurrentGinkgoTestDescription().TestText, "SD-CICD", "Michael Wilson", "sd-cicd-alerts", "sd-cicd@redhat.com", 4)
	})
	h := helper.New()

	prometheusTimeoutInSeconds := 900
	ginkgo.It("should include Prometheus data", func() {
		// setup runner
		// this command is has specific code to capture and suppress an exit code of
		// 1 as tar 1.26 will exit 1 if files change while the tar is running, as is
		// common for a running prometheus instance
		cmd := promCollectCmd + " >" + runner.DefaultRunner.OutputDir + "/prometheus.tar.gz\" ; err=$? ; if (( $err != 1 )) ; then exit $err ; fi"
		h.SetServiceAccount("system:serviceaccount:%s:cluster-admin")
		r := h.Runner(cmd)
		r.Name = "collect-prometheus"

		// run tests
		stopCh := make(chan struct{})
		err := r.Run(prometheusTimeoutInSeconds, stopCh)
		Expect(err).NotTo(HaveOccurred())

		// get results
		results, err := r.RetrieveResults()
		Expect(err).NotTo(HaveOccurred())

		// write results
		h.WriteResults(results)
	}, float64(prometheusTimeoutInSeconds+30))
})
