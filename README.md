# aws lambda golang client

## Example

- Run you golang lambda:

```bash
export _LAMBDA_SERVER_PORT=55555
./your_lambda
```

- Tringging your lambda function

```bash
export _LAMBDA_SERVER_PORT=55555
echo '{"foo": "bar"}' | ./labda_client
```
