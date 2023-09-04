type CreditCardProps = {
    card: {
        last4: string;
        brand: string;
        exp_month: number;
        exp_year: number;
    };
};

export default function CreditCard (props: CreditCardProps) {
    const { last4, brand, exp_month, exp_year } = props.card
    return (
        <option>
            {brand} **** **** **** {last4} - exp: {exp_month}/{exp_year}
        </option>
    )
}