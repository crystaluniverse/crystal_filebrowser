FROM onlyoffice/documentserver:5.6
ADD ./startup.sh /start_docuserver_grid.sh
RUN chmod +x /start_docuserver_grid.sh
ENTRYPOINT [ "/bin/sh", "-c", "/start_docuserver_grid.sh"]
