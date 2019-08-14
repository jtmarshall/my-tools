# SKUle

Single-page application to process sku's for the paid search team. Created with Python / React, and Zappa for Lambda deployment.

Hosted on AWS Lambda this repo contains both server-side and front-end source code.

- Uploads larger than a specified size are sent to Amazon S3 first because Lambda's maximum upload limit.

## Getting Started / Installing

Clone repo, navigate to 'static' folder for front-end, and run:
```
npm install
```

After packages are installed, run local dev server from 'server' folder:
```
python server.py
```

### Prerequisites

Things you need to have installed
```
python3
node
npm
```

## Deployment

Build front-end in 'static' folder:
```
npm run build
```

## Built With

* [React](https://reactjs.org/) - Library used for UI
* [Lambda](https://aws.amazon.com/lambda/) - Serverless data processing

## Authors

* **Jordan Marshall**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
