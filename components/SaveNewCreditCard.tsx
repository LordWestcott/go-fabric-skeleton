"use client";
import { CardElement, useElements, useStripe } from "@stripe/react-stripe-js";
import { SetupIntentResult, StripeCardElement } from "@stripe/stripe-js";

type Props = {
  setupIntent: any;
  setSetupIntent: (arg: any) => void;
  getWallet: () => void;
};

export default function SaveNewCreditCard(props: Props) {
  const elements = useElements();
  const stripe = useStripe();

  const handleSubmit = async (event: any) => {
    const cardElement = elements?.getElement(CardElement);

    const result = await stripe?.confirmCardSetup(
      props.setupIntent.client_secret,
      {
        payment_method: {
          card: cardElement as StripeCardElement,
        },
      }
    );

    if (result) {
      const { setupIntent: updatedSetupIntent, error } = result;

      if (error) {
        alert(error.message);
        console.log(error);
      } else {
        props.setSetupIntent(updatedSetupIntent);
        await props.getWallet();
        alert("Success Card Saved!");
      }
    } else {
      alert("Error saving card");
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <CardElement/>
      {/* <CardElement options={{hidePostalCode: true}}/> */}
      <button type="submit">Save Card</button>
    </form>
  );
}
