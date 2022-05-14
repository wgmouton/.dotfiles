FROM arch:latest

COPY ./ /.dotfiles
RUN chmod -R +x ./ && ./install.sh


ENTRYPOINT [ "./utilities/dev-environment/setup-dev-environment.sh", "local", "./workspace"]