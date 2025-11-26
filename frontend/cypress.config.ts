const { defineConfig } = require('cypress')

module.exports = defineConfig({
    e2e: {
        baseUrl: process.env.CYPRESS_baseUrl || 'http://localhost:3000',
        setupNodeEvents(on, config) {
            // implement node event listeners here
        },
        supportFile: 'cypress/support/e2e.js',
        specPattern: 'cypress/e2e/**/*.cy.{js,jsx,ts,tsx}',
        video: false,
        screenshotOnRunFailure: true,
    },
    env: {
        apiUrl: process.env.CYPRESS_apiUrl || 'http://localhost:8080/api',
    },
})
