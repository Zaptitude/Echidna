package wp

import (
	"context"
	"os"

	"github.com/Paraflare/Echidna/pkg/vulnerabilities"
)

// Scan will call the vulnerability packages scanning function to check each file for vulns
// if it finds vulns the plugin will be moved to the inspect/ folder with the results stored
// with it as a .txt file with the same name
func scanWorker(ctx context.Context, errCh chan error, workQueue chan *Plugin, resultsQueue chan vulnerabilities.Results) {

	for p := range workQueue {
		scanResults := vulnerabilities.Results{
			Plugin:  p.Name,
			Modules: make(map[string][]vulnerabilities.VulnResults),
		}

		err := vulnerabilities.ZipScan(ctx, p.OutPath, &scanResults)
		if err != nil {
			errCh <- err
			removeZip(p.OutPath, errCh)
			continue
		}
		if len(scanResults.Modules) > 0 {
			err := p.moveToInspect(&scanResults)
			if err != nil {
				errCh <- err
				removeZip(p.OutPath, errCh)
				continue
			}
			err = p.saveResults(&scanResults)
			if err != nil {
				errCh <- err
				continue
			}

			resultsQueue <- scanResults
		}
		removeZip(p.OutPath, errCh)
	}

}

func resultsWorker(ctx context.Context, errCh chan error, plugins *Plugins, queue chan vulnerabilities.Results) {

	for result := range queue {

		plugins.resMu.Lock()

		plugins.LatestVuln = result
		plugins.VulnsFound++
		plugins.Vulns = append(plugins.Vulns, result)
		plugins.FilesScanned++

		plugins.resMu.Unlock()
	}

}

func removeZip(path string, errCh chan error) {
	err := os.Remove(path)
	if err != nil {
		errCh <- err
	}
}