# Example Dockerfile snippet (using vsftpd)
FROM alpine:latest
RUN apk add --update vsftpd openssl
# Configure vsftpd for implicit FTPS on port 990
RUN echo "ssl_enable=YES" >> /etc/vsftpd/vsftpd.conf
RUN echo "implicit_ssl=YES" >> /etc/vsftpd/vsftpd.conf
RUN echo "listen_port=990" >> /etc/vsftpd/vsftpd.conf
# Generate self-signed certificate (for testing purposes only)
RUN openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/ssl/private/vsftpd.pem -out /etc/ssl/certs/vsftpd.pem -subj "/CN=localhost"
RUN echo "rsa_cert_file=/etc/ssl/certs/vsftpd.pem" >> /etc/vsftpd/vsftpd.conf
RUN echo "rsa_private_key_file=/etc/ssl/private/vsftpd.pem" >> /etc/vsftpd/vsftpd.conf
EXPOSE 990
CMD ["vsftpd", "/etc/vsftpd/vsftpd.conf"]