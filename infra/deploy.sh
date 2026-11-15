aws cloudformation deploy \
  --template-file cloudformation.yml \
  --stack-name ns-scoring-terraform-backend \
  --capabilities CAPABILITY_NAMED_IAM
