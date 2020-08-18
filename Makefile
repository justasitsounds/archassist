start-puml:
	docker run -d -v $(AWSICONDIST):/include -e ALLOW_PLANTUML_INCLUDE=true -p 8080:8080 plantuml/plantuml-server:tomcat