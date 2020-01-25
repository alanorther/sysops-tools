key=$1
secret=$2

# Pull down container
aws s3 cp s3://verve-yext-integration-container/portal.tar.Z ./

# Load container locally
docker load < portal.tar.z
# Output
# Find container image id
image_id=`docker images | grep portal | awk '{ print $3 }'`

# Tag image
echo "docker tag $image_id 102126644248.dkr.ecr.us-west-2.amazonaws.com/yext-integration-server:master-latest"
docker tag $image_id 102126644248.dkr.ecr.us-west-2.amazonaws.com/yext-integration-server:master-latest

# Login to ecr
lgin=`aws ecr get-login --no-include-email --region us-west-2`

# Copy and paste that command to actually log in to the the registry
echo $lgin
$lgin

# Push your local image to the registry
docker images
echo "docker push 102126644248.dkr.ecr.us-west-2.amazonaws.com/yext-integration-server:master-latest"
docker push 102126644248.dkr.ecr.us-west-2.amazonaws.com/yext-integration-server:master-latest

# If thats succesful you should be able to deploy
echo "docker run -e AWS_ACCESS_KEY_ID=$key -e AWS_SECRET_ACCESS_KEY=$secret -e VERSION=latest silintl/ecs-deploy:3.4.0 -t 600 -r us-west-2 -c qa -n yext-integration-server -i 102126644248.dkr.ecr.us-west-2.amazonaws.com/yext-integration-server:master-latest --min 50 --max 200 --enable-rollback"
docker run -e AWS_ACCESS_KEY_ID=$key -e AWS_SECRET_ACCESS_KEY=$secret -e VERSION=latest silintl/ecs-deploy:3.4.0 -t 600 -r us-west-2 -c qa -n yext-integration-server -i 102126644248.dkr.ecr.us-west-2.amazonaws.com/yext-integration-server:master-latest --min 50 --max 200 --enable-rollback
echo "docker run -e AWS_ACCESS_KEY_ID=$key -e AWS_SECRET_ACCESS_KEY=$secret -e VERSION=latest silintl/ecs-deploy:3.4.0 -t 600 -r us-west-2 -c staging -n yext-integration-server -i 102126644248.dkr.ecr.us-west-2.amazonaws.com/yext-integration-server:master-latest --min 50 --max 200 --enable-rollback"
docker run -e AWS_ACCESS_KEY_ID=$key -e AWS_SECRET_ACCESS_KEY=$secret -e VERSION=latest silintl/ecs-deploy:3.4.0 -t 600 -r us-west-2 -c staging -n yext-integration-server -i 102126644248.dkr.ecr.us-west-2.amazonaws.com/yext-integration-server:master-latest --min 50 --max 200 --enable-rollback
#echo "docker run -e AWS_ACCESS_KEY_ID=$key -e AWS_SECRET_ACCESS_KEY=$secret -e VERSION=latest silintl/ecs-deploy:3.4.0 -t 600 -r us-east-1 -c production -n yext-integration-server -i 102126644248.dkr.ecr.us-west-2.amazonaws.com/yext-integration-server:master-latest --min 50 --max 200 --enable-rollback"
#docker run -e AWS_ACCESS_KEY_ID=$key -e AWS_SECRET_ACCESS_KEY=$secret -e VERSION=latest silintl/ecs-deploy:3.4.0 -t 600 -r us-east-1 -c production -n yext-integration-server -i 102126644248.dkr.ecr.us-west-2.amazonaws.com/yext-integration-server:master-latest --min 50 --max 200 --enable-rollback

#clean up images
docker rmi -f $(docker images -q)
