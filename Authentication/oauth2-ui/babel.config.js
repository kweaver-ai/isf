module.exports = {
    env: {
        development: {
            presets: [["next/babel", { "preset-env": { useBuiltIns: "entry", corejs: 3 } }]],
        },
        production: {
            presets: [["next/babel", { "preset-env": { useBuiltIns: "entry", corejs: 3 } }]],
        },
        test: {
            presets: [["next/babel", { "preset-env": { modules: "commonjs" } }]],
        },
    },
};
