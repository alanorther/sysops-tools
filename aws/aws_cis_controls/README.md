# AWS CIS Scripts

A collection of AWS scripts to test the CIS recommendations.

### Prerequisites

You need Python and boto3 python library.

### Installing

```
pip install boto3
```

## Running the tests

You can run the scripts by environment and region, otherwise will use your default. 

```
python  aws_cis_mfa_use.py --env staging --region us-east-1
```

