
FROM alpine as build

ARG BINNAME
WORKDIR /
COPY ${BINNAME} /
RUN mv /${BINNAME} /app

FROM gcr.io/distroless/static:nonroot

USER 65532:65532
WORKDIR /
COPY --from=build /app /app
ENTRYPOINT ["/app"]
