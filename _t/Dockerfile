FROM alpine:latest
RUN apk --update add openjdk8-jre curl screen && \
    curl -s -L https://marketplace.atlassian.com/download/plugins/atlassian-plugin-sdk-tgz | tar xzf - && \
    ln -s /atlassian* /atlassian

ENV PATH=/bin:/usr/bin:/atlassian/bin

# Copy in the serivce and also the root .m2 settings to force cache everything.
# We also copy in /root/.java settings to prevent the dumb spam prompt from
# the atlas-run command:
# Would you like to subscribe to the Atlassian developer mailing list? (Y/y/N/n) Y: :
COPY dockerroot /
WORKDIR /jiratestservice

EXPOSE 8080

# we wrap the command with screen so that the dumb atlas-run has a tty to watch. Without screen
# there is no tty so atlas-run will immediately read an EOF (aka CTRL-D) and interpret that to
# mean we want the service to begin the "graceful shutdown" and exit
CMD ["screen", "-DmL", "atlas-run", "--http-port", "8080", "--context-path", "ROOT", "--server", "localhost"]
