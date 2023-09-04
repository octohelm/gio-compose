FROM golang:1.21-bullseye

ENV DEBIAN_FRONTEND=noninteractive

RUN \
    apt-get update \
    && apt-get install -y \
      gcc \
      pkg-config \
      libwayland-dev \
      libx11-dev \
      libx11-xcb-dev \
      libxkbcommon-x11-dev \
      libgles2-mesa-dev \
      libegl1-mesa-dev \
      libffi-dev \
      libxcursor-dev \
      libvulkan-dev \
      make

RUN \
    go install gioui.org/cmd/gogio@latest
