package src

type Pipeline struct {
	config  Config
	emailer Emailer
	log     Logger
}

func (p *Pipeline) run(project Project) {
	testsPassed := p.runTests(project)
	deploySuccessful := p.deploy(project, testsPassed)

	p.sendEmailSummary(testsPassed, deploySuccessful)
}

func (p *Pipeline) sendEmailSummary(testsPassed bool, deploySuccessful bool) {
	if !p.config.sendEmailSummary() {
		p.log.info("Email disabled")
		return
	}

	p.log.info("Sending email")

	if !testsPassed {
		p.emailer.send("Tests failed")
		return
	}

	if !deploySuccessful {
		p.emailer.send("Deployment failed")
		return
	}

	p.emailer.send("Deployment completed successfully")
}

func (p *Pipeline) deploy(project Project, testsPassed bool) (deploySuccessful bool) {
	if !testsPassed {
		return false
	}

	if project.deploy() == "success" {
		p.log.info("Deployment successful")
		return true
	} else {
		p.log.error("Deployment failed")
		return false
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
