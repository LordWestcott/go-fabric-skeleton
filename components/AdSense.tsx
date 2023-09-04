"use client"
import { useEffect } from "react";

declare const window: any;

const AdSense = ({ adSlot } : { adSlot: string; }) => {
	useEffect(() => {
		if (window) {
			(window.adsbygoogle = window.adsbygoogle || []).push({});
		}
	}, []);

	return(
	<ins
		className="adsbygoogle"
		style={{ display: "block" }}
		data-ad-client="pub-5931480780360258"
		data-ad-slot={adSlot}
		data-ad-format="auto"
		data-full-width-responsive="true"
	>	
	
	</ins>
	)
};

export default AdSense;