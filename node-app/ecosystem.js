{"apps": [
    {
      name: "Bench",
      script: "index.js",
      instances: "max",
      exec_mode: "cluster",
      env: {
        PORT: "8080",
      },
    },
  ],
};
