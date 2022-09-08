FROM migrate/migrate:v4.15.1
COPY migrations /migrations
ARG POSTGRES_USER
ARG POSTGRES_PASSWORD
ARG POSTGRES_HOST
ARG POSTGRES_DATABASE
ENTRYPOINT []
CMD ["sh", "-c", "migrate -path=/migrations/ -database \"postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_DATABASE}?sslmode=disable\" up"]