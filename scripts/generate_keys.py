# Run at the account root!

import os

print("Creating an rsa256 key pair")

generate_dir = "./keys"
if not os.path.exists(generate_dir):
    os.makedirs(generate_dir)
print("Key's name:")
name = input()

_ = os.popen(
    "openssl genpkey -algorithm RSA -out %s/rsa_private_%s.pem -pkeyopt rsa_keygen_bits:2048"
    % (generate_dir, name)
).read()

_ = os.popen(
    "openssl rsa -in %s/rsa_private_%s.pem -pubout -out %s/rsa_public_%s.pem"
    % (generate_dir, name, generate_dir, name)
).read()
