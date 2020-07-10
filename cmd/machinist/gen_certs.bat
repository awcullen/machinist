mkdir .\pki

(
  echo [req]
  echo distinguished_name=req
  echo [san]
  echo subjectAltName=URI:urn:127.0.0.1:device2,IP:127.0.0.1
) > .\pki\device2.cnf

openssl req -x509 -newkey rsa:2048 -sha256 -days 3650 -nodes ^
 -keyout .\pki\device2.key -out .\pki\device2.crt ^
 -subj /CN=device2 ^
 -extensions san ^
 -config .\pki\device2.cnf
