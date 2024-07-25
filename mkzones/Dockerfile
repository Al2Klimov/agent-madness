FROM golang as build

ADD . /mkzones
WORKDIR /mkzones

ENV CGO_ENABLED 0
RUN ["go", "build", "-ldflags", "-s -w", "."]


FROM scratch

COPY --from=build /mkzones/mkzones /
USER 101:101
CMD ["/mkzones"]
