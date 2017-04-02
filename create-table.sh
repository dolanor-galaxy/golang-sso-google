aws dynamodb --profile default --region us-west-2 create-table \
    --table-name Users \
    --key-schema \
        AttributeName=email,KeyType=HASH \
    --attribute-definitions \
        AttributeName=email,AttributeType=S \
    --provisioned-throughput \
        ReadCapacityUnits=1,WriteCapacityUnits=1
