import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

interface IndividualParam {
    individualID: string;
}

interface Association {
    id: string;
    description: string;
}

interface IndividualResult {
    id: string;
    firstName: string;
    lastName: string;
    role: string;
    title: string;
    twitter: string;
    associations: Association[];
}

interface Contribution {
    id: string;
    amount: number;
    date: Date;
    contributorName: string;
    contributorId: string;
    recipientName: string;
    recipientId: string;
}

const Individual = () => {
    const { individualID } = useParams<IndividualParam>();
    const [individual, setIndividual] = useState<IndividualResult>();
    const [contribReceived, setContribReceived] = useState<Contribution[]>();
    const [contribGiven, setContribGiven] = useState<Contribution[]>();

    useEffect(() => {
        // Get individual
        fetch("/get-individual", {
            method: "POST",
            body: JSON.stringify({
                id: individualID,
            }),
        })
            .then((resp) => resp.json())
            .then((data) => {
                setIndividual({
                    id: data.id,
                    firstName: data.first_name,
                    lastName: data.last_name,
                    role: data.role,
                    title: data.title,
                    twitter: data.twitter,
                    associations: data.associations?.map((a: any) => {
                        return {
                            id: a.id,
                            description: a.description,
                        };
                    }),
                });
            })
            .catch((err) => {
                console.log("error: ", err);
            });

        // Get contributions received
        fetch("/individual-contributions-received", {
            method: "POST",
            body: JSON.stringify({
                id: individualID,
            }),
        })
            .then((resp) => resp.json())
            .then((data) => {
                setContribReceived(
                    data?.map((d: any) => {
                        return {
                            id: d.id,
                            amount: d.amount,
                            date: d.date,
                            contributorName: d.contributor_name,
                            contributorId: d.contributor_id,
                            recipientName: d.recipient_name,
                            recipientId: d.recipient_id,
                        };
                    })
                );
            })
            .catch((err) => {
                console.log("error: ", err);
            });

        // Get contributions given
        fetch("/individual-contributions-given", {
            method: "POST",
            body: JSON.stringify({
                id: individualID,
            }),
        })
            .then((resp) => resp.json())
            .then((data) => {
                setContribGiven(
                    data?.map((d: any) => {
                        return {
                            id: d.id,
                            amount: d.amount,
                            date: d.date,
                            contributorName: d.contributor_name,
                            contributorId: d.contributor_id,
                            recipientName: d.recipient_name,
                            recipientId: d.recipient_id,
                        };
                    })
                );
            })
            .catch((err) => {
                console.log("error: ", err);
            });
    }, []);

    return individual ? (
        <>
            <h3>{individual.firstName + " " + individual.lastName}</h3>
            <p>{individual.role}</p>
            {individual.associations?.map((a: Association) => {
                return <p key={a.id}>{a.description}</p>;
            })}
            {contribGiven && (
                <>
                    <h3>Contributions Given</h3>
                    {contribGiven?.map((c: Contribution) => {
                        return (
                            <p key={c.id}>
                                {c.contributorName + ", " + c.amount}
                            </p>
                        );
                    })}
                </>
            )}
            {contribReceived && (
                <>
                    <h3>Contributions Received</h3>
                    {contribReceived?.map((c: Contribution) => {
                        return (
                            <p key={c.id}>
                                {c.contributorName + ", " + c.amount}
                            </p>
                        );
                    })}
                </>
            )}
        </>
    ) : (
        <></>
    );
};

export default Individual;
