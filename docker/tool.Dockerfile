# azure-functions-core-tools only available in amd64
# issues at https://github.com/Azure/azure-functions-core-tools/issues/2901
FROM --platform=linux/amd64 debian:bullseye

RUN apt update
RUN apt install curl unzip -y ca-certificates apt-transport-https lsb-release gnupg

# azure-functions-core-tools need libicu for ICU https://github.com/dotnet/core/issues/2186
RUN apt install libicu-dev -y

# Install Azure Functions Core tools
RUN curl -L 'https://github.com/Azure/azure-functions-core-tools/releases/download/4.0.4483/Azure.Functions.Cli.linux-x64.4.0.4483.zip' -o Azure.Functions.Cli.zip && \
    mkdir -p /opt/func && \
    unzip -d /opt/func Azure.Functions.Cli.zip && \
    chmod +x /opt/func/func && \
    chmod +x /opt/func/gozip && \
    rm Azure.Functions.Cli.zip

# Add func to path
ENV PATH="/opt/func:${PATH}"

# Install AZ CLI
RUN mkdir -p /etc/apt/keyrings
RUN curl -sLS https://packages.microsoft.com/keys/microsoft.asc | \
        gpg --dearmor | \
        tee /etc/apt/keyrings/microsoft.gpg > /dev/null
RUN chmod go+r /etc/apt/keyrings/microsoft.gpg

RUN echo "deb [arch=`dpkg --print-architecture` signed-by=/etc/apt/keyrings/microsoft.gpg] https://packages.microsoft.com/repos/azure-cli/ $(lsb_release -cs) main" | \
    tee /etc/apt/sources.list.d/azure-cli.list

RUN apt-get install azure-cli -y

WORKDIR /api