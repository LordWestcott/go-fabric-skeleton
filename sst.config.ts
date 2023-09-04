import { SSTConfig } from "sst";
import { Api, Function, NextjsSite } from "sst/constructs";

export default {
  config(_input) {
    return {
      name: "appointme",
      region: "eu-north-1",
    };
  },
  stacks(app) {

    app.setDefaultFunctionProps({
      runtime: "go1.x",
    });

    app.stack(function Site({ stack }) {
      require('dotenv').config()

      const func_Test = new Function(stack, "func-test", {
        handler: "functions/lambda/test/main.go",
        runtime: "go1.x",
        environment: {
        },
      });

      const api = new Api(stack, "gofabric-api", {
        routes: {
          "GET /": func_Test,
        }
      });

      let base_environment = {
        STRIPE_PUBLISHABLE_KEY: process.env.STRIPE_PUBLISHABLE_KEY || "",
        STRIPE_PRIVATE_KEY: process.env.STRIPE_PRIVATE_KEY || "",
        STRIPE_WEBHOOK_SECRET: process.env.STRIPE_WEBHOOK_SECRET || "",
        TWILIO_ACCOUNT_SID: process.env.TWILIO_ACCOUNT_SID || "",
        TWILIO_AUTH_TOKEN: process.env.TWILIO_AUTH_TOKEN || "",
        TWILIO_VERIFY_SERVICE_ID: process.env.TWILIO_VERIFY_SERVICE_ID || "",
        TWILIO_WHATSAPP_NUMBER: process.env.TWILIO_WHATSAPP_NUMBER || "",
        SENDGRID_API_KEY: process.env.SENDGRID_API_KEY || "",
        SMS_SERVICE: process.env.SMS_SERVICE || "",
        EMAIL_SERVICE: process.env.EMAIL_SERVICE || "",
        WHATSAPP_SERVICE: process.env.WHATSAPP_SERVICE || "",
        VERIFICATION_SERVICE: process.env.VERIFICATION_SERVICE || "",
        OPENAI_API_KEY: process.env.OPENAI_API_KEY || "",
        DATABASE_URL: process.env.DATABASE_URL || "",
        SIGNING_KEY: process.env.SIGNING_KEY || "",
        GOOGLE_CLIENT_ID: process.env.GOOGLE_CLIENT_ID || "",
        GOOGLE_CLIENT_SECRET: process.env.GOOGLE_CLIENT_SECRET || "",
        GOOGLE_STATE: process.env.GOOGLE_STATE || "",
        GOOGLE_REDIRECT: process.env.GOOGLE_REDIRECT || "",
      
        REDIRECT_SOMETHING_WENT_WRONG_PAGE: process.env.REDIRECT_SOMETHING_WENT_WRONG_PAGE || "",
        REDIRECT_LOGIN_SUCCESS: process.env.REDIRECT_LOGIN_SUCCESS || "",

        APP_NAME: process.env.APP_NAME || "",
        SESSION_TYPE: process.env.SESSION_TYPE || "",

        COOKIE_NAME: process.env.COOKIE_NAME || "",
        COOKIE_LIFETIME: process.env.COOKIE_LIFETIME || "",
        COOKIE_PERSISTS: process.env.COOKIE_PERSISTS || "",
        COOKIE_DOMAIN: process.env.COOKIE_DOMAIN || "",
        COOKIE_SECURE: process.env.COOKIE_SECURE || "",

        FRONTEND_URL: process.env.FRONTEND_URL || "",
        HOST: api.url|| "",

        JWT_SECRET: process.env.JWT_SECRET || "",

        GOOGLE_ADSENSE_PUBLISHER_ID: process.env.GOOGLE_ADSENSE_PUBLISHER_ID || "",
        
      }


      console.log(base_environment)

      const site = new NextjsSite(stack, "site", {
        environment: {
          NEXT_PUBLIC_API_URL: base_environment.HOST,
          NEXT_PUBLIC_GOOGLE_ADSENSE_PUBLISHER_ID: base_environment.GOOGLE_ADSENSE_PUBLISHER_ID,
          NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY: base_environment.STRIPE_PUBLISHABLE_KEY,
          NEXT_PUBLIC_GA_MEASUREMENT_ID: process.env.GA_MEASUREMENT_ID || "",
        }
      });
      const siteurl = site.url || process.env.FRONTEND_URL || "";
      base_environment = {
        ...base_environment,
        FRONTEND_URL: siteurl,
      }


      //SECTION - AUTH
      const func_GoogleSignin = new Function(stack, "func-google-signin", {
        handler: "functions/lambda/auth/google-signin/main.go",
        runtime: "go1.x",
        environment: {
          ...base_environment,
        },
      });

      const func_GoogleSigninCallback = new Function(stack, "func-google-signin-callback", {
        handler: "functions/lambda/auth/google-signin-callback/main.go",
        runtime: "go1.x",
        environment: {
          ...base_environment,
        },
      });

      const func_Session = new Function(stack, "func-session", {
        handler: "functions/lambda/auth/session/main.go",
        runtime: "go1.x",
        environment: {
          ...base_environment
        },
      });

      api.addRoutes(stack, {
          "GET /google-signin": func_GoogleSignin,
          "GET /google-signin/callback": func_GoogleSigninCallback,
          "GET /session": func_Session,      
        })

      //SECTION - Stripe Examples
      const func_Stripe_CreateExternalCheckoutSession = new Function(stack, "func-stripe-create-external-checkout-session", {
        handler: "functions/lambda/example/stripe/create-external-checkout/main.go",
        runtime: "go1.x",
        environment: {
          ...base_environment,
        },
      });

      const func_Stripe_CreateAnonPaymentIntent = new Function(stack, "func-stripe-create-payment-intent", {
        handler: "functions/lambda/example/stripe/create-anon-payment-intent/main.go",
        runtime: "go1.x",
        environment: {
          ...base_environment,
        },
      });

      const func_Stripe_CreateSetupIntent = new Function(stack, "func-stripe-create-setup-intent", {
        handler: "functions/lambda/example/stripe/create-setup-intent/main.go",
        runtime: "go1.x",
        environment: {
          ...base_environment,
        },
      });

      const func_Stripe_ListPaymentMethods = new Function(stack, "func-stripe-list-payment-methods", {
        handler: "functions/lambda/example/stripe/list-payment-methods/main.go",
        runtime: "go1.x",
        environment: {
          ...base_environment,
        },
      });

      const func_Stripe_Webhook = new Function(stack, "func-stripe-webhook", {
        handler: "functions/lambda/example/stripe/webhook/main.go",
        runtime: "go1.x",
        environment: {
          ...base_environment,
        },
      });

      const func_Stripe_CreateSubscription = new Function(stack, "func-stripe-create-subscription", {
        handler: "functions/lambda/example/stripe/create-subscription/main.go",
        runtime: "go1.x",
        environment: {
          ...base_environment,
        },
      });

      const func_Stripe_CancelSubscription = new Function(stack, "func-stripe-cancel-subscription", {
        handler: "functions/lambda/example/stripe/cancel-subscription/main.go",
        runtime: "go1.x",
        environment: {
          ...base_environment,
        },
      });

      const func_Stripe_ListUserSubscriptions = new Function(stack, "func-stripe-list-user-subscriptions", {
        handler: "functions/lambda/example/stripe/list-users-subscriptions/main.go",
        runtime: "go1.x",
        environment: {
          ...base_environment,
        },
      });

      api.addRoutes(stack, {
          "POST /stripe/create-external-checkout-session": func_Stripe_CreateExternalCheckoutSession,
          "POST /stripe/create-anon-payment-intent": func_Stripe_CreateAnonPaymentIntent,
          "POST /stripe/webhook": func_Stripe_Webhook,
          "POST /stripe/wallet": func_Stripe_CreateSetupIntent,
          "GET /stripe/wallet": func_Stripe_ListPaymentMethods,
          "GET /stripe/subscriptions": func_Stripe_ListUserSubscriptions,
          "POST /stripe/create-subscription": func_Stripe_CreateSubscription,
          "PATCH /stripe/cancel-subscription": func_Stripe_CancelSubscription,
      });


      //Products
      const func_Products_ListProducts = new Function(stack, "func-products-list-products", {
        handler: "functions/lambda/products/list-products/main.go",
        runtime: "go1.x",
        environment: {
          ...base_environment,
        },
      });

      api.addRoutes(stack, {
          "GET /products": func_Products_ListProducts,
      });
        

      stack.addOutputs({
        SiteUrl: site.url,
        ApiEndpoint: api.url,
      });
    });
  },
} satisfies SSTConfig;
