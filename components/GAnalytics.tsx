import Script from "next/script";

//this function hooks up gtag.js to the page

type GAnalyticsProps = {
    GA_MEASUREMENT_ID: string;
}

export default function GAnalytics({ GA_MEASUREMENT_ID }: GAnalyticsProps) {

  return (
    <>
      <Script src={`https://www.googletagmanager.com/gtag/js?id=${GA_MEASUREMENT_ID}`} />
      <Script id="google-analytics">
        {`
                    window.dataLayer = window.dataLayer || [];
                    function gtag(){dataLayer.push(arguments);}
                    gtag('js', new Date());
            
                    gtag('config', '` + GA_MEASUREMENT_ID + `');
                    `}
      </Script>
    </>
  );
}
