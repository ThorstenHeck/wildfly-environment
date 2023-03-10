# Deployment Zip

The input data for the deployment is a ZIP File with following defined structure:

.zip
- jcli (scripts to execute on the Wildfly/EAP server - jndi names, user management etc.)
- sql ()
- modules ()
- app (.ear and .war files to deploy)
- 


## example zip

mkdir deployment_example
mkdir deployment_example/jcli
mkdir deployment_example/sql
mkdir deployment_example/modules
mkdir deployment_example/app


zip -r deployment_example.zip deployment_example

## upload to operator

curl -X POST http://localhost:10000/deploy

curl -F file=@deployment_example.zip http://localhost:10000/upload


