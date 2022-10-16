package src

type Pipeline struct {
	config  Config
	emailer Emailer
	log     Logger
}

func (p *Pipeline) run(project Project) {
	var deploySuccessful bool

	testsPassed := p.runTests(project)

	if testsPassed {
		if "success" == project.deploy() {
			p.log.info("Deployment successful")
			deploySuccessful = true
		} else {
			p.log.error("Deployment failed")
			deploySuccessful = false
		}
	} else {
		deploySuccessful = false
	}

	if p.config.sendEmailSummary() {
		p.log.info("Sending email")
		if testsPassed {
			if deploySuccessful {
				p.emailer.send("Deployment completed successfully")
			} else {
				p.emailer.send("Deployment failed")
			}
		} else {
			p.emailer.send("Tests failed")
		}
	} else {
		p.log.info("Email disabled")
	}
}

func (p *Pipeline) runTests(project Project) (testsPassed bool) {
	if !project.hasTests() {
		p.log.info("No tests")
		return true
	}

	if project.runTests() == "success" {
		p.log.info("Tests passed")
		return true
	} else {
		p.log.error("Tests failed")
		return false
	}
}
