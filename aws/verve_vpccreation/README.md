laptop ~/puppet/verve/modules/verve_vpccreation/files$ ./createvpc.sh onetest 10.160.0.0 stage
Would you like to auto deploy this stack? [y/n]y
Copying cloud formation config to s3...
../templates/vpc-onetest-stage-template.json -> s3://verve-cf-templates/stage/vpc-onetest-stage-template.json  [1 of 1]
aws --region us-west-2 cloudformation create-stack --stack-name onetest --template-url https://s3.amazonaws.com/verve-cf-templates/stage/vpc-onetest-stage-template.json
Creating machine configs...

Once the stack is done building, You can deploy all nodes in the environment with this command:
for line in `ls verve/modules/verve_aws/manifests/onetest/`;do puppet apply --hiera_config hiera.yaml --modulepath=site-modules:verve/modules verve/modules/verve_aws/manifests/onetest/$line -t --noop;done
Or you can deploy the nodes one at a tim

And dont be fooled by the "Automatic" deploy question, it doesnt auto deploy yet it just spits out the command you want to use to deploy the cloudform from your laptop which would look like this:
laptop ~/puppet$ aws --region us-west-2 cloudformation create-stack --stack-name onetest --template-url https://s3.amazonaws.com/verve-cf-templates/stage/vpc-onetest-stage-template.json
{ "StackId": "arn:aws:cloudformation:us-west-2:102126644248:stack/onetest/777a47f0-aa7b-11e5-bfce-50e2414b0a7c" }
