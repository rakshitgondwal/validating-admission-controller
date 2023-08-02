FROM alpine

COPY validating-admission-controller /validating-admission-controller 

ENTRYPOINT [ "./validating-admission-controller" ]