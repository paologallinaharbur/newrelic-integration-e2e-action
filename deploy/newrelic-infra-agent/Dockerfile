FROM newrelic/infrastructure:latest

RUN mkdir -p /usr/local/prometheus-exporters/bin
RUN mkdir -p /var/db/newrelic-infra/custom-integrations/bin

ADD integrations.d/nri-config.yml /etc/newrelic-infra/integrations.d
ADD bin/* /var/db/newrelic-infra/newrelic-integrations/bin
ADD exporters/* /usr/local/prometheus-exporters/bin

RUN chmod +x /var/db/newrelic-infra/newrelic-integrations/bin/*
RUN chmod +x /usr/local/prometheus-exporters/bin/*

#RUN cp /var/db/newrelic-infra/newrelic-integrations/bin/nri-prometheus /var/db/newrelic-infra/custom-integrations/bin/nri-prometheus
#RUN chmod +x /var/db/newrelic-infra/custom-integrations/bin/*

RUN rm /etc/newrelic-infra/integrations.d/docker-config.yml
