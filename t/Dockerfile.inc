RUN curl -Ls https://raw.githubusercontent.com/cptactionhank/docker-atlassian-jira-software/master/docker-entrypoint.sh -o /docker-entrypoint.sh \
    && chmod 755 /docker-entrypoint.sh \
    && echo jira.websudo.is.disabled = true >> /var/atlassian/jira/jira-config.properties