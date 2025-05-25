window.onload = function() {
  const origin = window.location.origin;
  const specUrl = `${origin}/openapi.json`;
  const style = document.createElement("style");
  style.innerHTML = `
    .swagger-ui .topbar { display: none !important; }
    .swagger-ui .info { margin: 0 !important; }
  `;
  document.head.appendChild(style);
  window.ui = SwaggerUIBundle({
    url: specUrl,
    dom_id: "#swagger-ui",
    deepLinking: true,
    layout: "BaseLayout",
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl,
      () => ({
        wrapComponents: {
          Topbar: () => () => null,
        },
      }),
    ],
    supportedSubmitMethods: ["get", "post", "put", "delete", "patch"],
    docExpansion: "none",
    filter: false,
  });
};
