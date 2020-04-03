package common

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Analyzer -
type Analyzer struct {
	config  *BaseConfig
	metrics *Metrics
}

// NewAnalyzer -
func NewAnalyzer(config *BaseConfig, metrics *Metrics) *Analyzer {
	return &Analyzer{
		config:  config,
		metrics: metrics,
	}
}

// RunForever - starts endless analyze loop
func (a *Analyzer) RunForever(interval time.Duration) {
	go func() {
		for {
			log.Infof("running full analyze")
			start := time.Now()
			for _, cProject := range a.config.Projects {
				a.ProcessProject(&cProject)
			}
			a.metrics.Duration.Set(time.Now().Sub(start).Seconds())
			time.Sleep(interval)
		}
	}()
}

// ProcessProject - analyze a single project
func (a *Analyzer) ProcessProject(config *GitConfig) error {
	modules, err := a.analyzeProject(config)
	if err != nil {
		a.metrics.Status.WithLabelValues(config.URL).Set(float64(0))
		return err
	}
	main, dependencies := a.extractModule(modules)
	a.writeMetrics(config, main, dependencies)
	a.metrics.Status.WithLabelValues(config.URL).Set(float64(1))
	return nil
}

func (a *Analyzer) writeMetrics(config *GitConfig, main ModulePublic, deps []ModulePublic) {
	config.Entry().Debug("writing statistics")
	a.metrics.Info.WithLabelValues(main.Path, main.GoVersion).Set(1)
	for _, cDep := range deps {
		mValue := float64(0)
		mLatestVersion := cDep.Version
		mType := "direct"
		if cDep.Indirect {
			mType = "indirect"
		}
		if cDep.Update != nil {
			mValue = 1000.0
			mLatestVersion = cDep.Update.Version
			if cDep.Time != nil && cDep.Update.Time != nil {
				mValue = cDep.Update.Time.Sub(*cDep.Time).Hours() / 24.0
			}
		}
		a.metrics.Deprecated.WithLabelValues(
			main.Path, cDep.Path, mType,
			cDep.Version, mLatestVersion,
		).Set(mValue)
	}
}

func (a *Analyzer) getRepository(config *GitConfig, dir string) error {
	config.Entry().Debug("cloning repository")

	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               config.URL,
		SingleBranch:      true,
		Depth:             1,
		Auth:              config.AuthMethod(),
		RecurseSubmodules: git.NoRecurseSubmodules,
	})
	if err != nil {
		err = errors.Wrapf(err, "unable to checkout")
		config.Entry().Errorf(err.Error())
		return err
	}

	return nil
}

func (a *Analyzer) getModules(config *GitConfig, dir string) ([]ModulePublic, error) {
	config.Entry().Debug("extracting go modules")

	cmd := exec.Command("go", "list", "-u", "-m", "-json", "all")
	cmd.Dir = dir
	content, err := cmd.Output()
	if err != nil {
		if exerr, ok := err.(*exec.ExitError); ok {
			config.Entry().Errorf(string(exerr.Stderr))
		}
		err = errors.Wrap(err, "unable to run go analysis")
		config.Entry().Errorf(err.Error())
		return nil, err
	}

	jsonStr := string(content)
	jsonStr = strings.ReplaceAll(jsonStr, "}\n{", "},\n{")
	jsonStr = fmt.Sprintf("[%s]", jsonStr)
	modules := []ModulePublic{}
	err = json.Unmarshal([]byte(jsonStr), &modules)
	if err != nil {
		err = errors.Wrap(err, "unable to parse go list output")
		config.Entry().Errorf(err.Error())
		return nil, err
	}

	return modules, nil
}

func (a *Analyzer) extractModule(mods []ModulePublic) (ModulePublic, []ModulePublic) {
	var main ModulePublic
	res := []ModulePublic{}
	for _, cModule := range mods {
		if cModule.Main {
			main = cModule
		} else {
			res = append(res, cModule)
		}
	}
	return main, res
}

func (a *Analyzer) analyzeProject(config *GitConfig) ([]ModulePublic, error) {
	config.Entry().Info("analysing project")

	dir, err := ioutil.TempDir("", "git-checkout")
	if err != nil {
		err = errors.Wrap(err, "unable to create temp directory")
		config.Entry().Errorf(err.Error())
		return nil, err
	}
	defer os.RemoveAll(dir)

	if err = a.getRepository(config, dir); err != nil {
		return nil, err
	}

	return a.getModules(config, dir)
}

// Local Variables:
// ispell-local-dictionary: "american"
// End:
