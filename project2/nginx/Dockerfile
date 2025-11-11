FROM nginx:latest


WORKDIR /proxy

COPY generate-ssl.sh .

RUN chmod +x generate-ssl.sh && ./generate-ssl.sh 
RUN mkdir -p /etc/nginx/ssl/ &&  mv *.pem /etc/nginx/ssl/ 

COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
EXPOSE 443

CMD ["nginx", "-g", "daemon off;"]