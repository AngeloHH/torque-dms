# CAMBIO IMPORTANTE: Subimos a la versión 1.24 para arreglar el error de dependencias
FROM golang:1.24-alpine

# Instalamos lo básico: Git (necesario para Go), Node y NPM.
RUN apk add --no-cache git nodejs npm

# Definimos la carpeta de trabajo
WORKDIR /app

# NOTA: No copiamos nada (ni go.mod, ni package.json).
# En Codespaces, tus archivos de la carpeta actual aparecerán 
# automáticamente en /workspaces/tu-repositorio.

# Mantenemos el contenedor vivo
CMD ["tail", "-f", "/dev/null"]
