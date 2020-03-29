FROM scratch
COPY dist/backend /app
CMD ["./app"]