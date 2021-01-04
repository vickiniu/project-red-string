import React, { useState } from "react";
import Form from "react-bootstrap/Form";
import { Link } from "react-router-dom";

interface ResultsListProps {
    data: SearchResult[] | undefined;
}

interface SearchResult {
    id: string;
    firstName: string;
    lastName: string;
}

const ResultsList = (props: ResultsListProps) => {
    const { data } = props;
    return data ? (
        <ul>
            {data.map((d: SearchResult) => {
                return (
                    <li key={d.id}>
                        <Link to={`individual/${d.id}`}>
                            {d.firstName + " " + d.lastName}
                        </Link>
                    </li>
                );
            })}
        </ul>
    ) : (
        <></>
    );
};

const Search = () => {
    const [results, setResults] = useState<SearchResult[]>();

    function handleChange(event: React.ChangeEvent<HTMLInputElement>) {
        console.log(event.target.value);
        // TODO(vicki): put in env variable for host
        fetch("http://localhost:8080/search-individuals", {
            method: "POST",
            body: JSON.stringify({
                query: event.target.value,
            }),
        })
            .then((resp) => resp.json())
            .then((data) => {
                console.log("data: ", data);
                let transformedResults: SearchResult[] = data?.map((d: any) => {
                    return {
                        id: d?.id,
                        firstName: d?.first_name,
                        lastName: d?.last_name,
                    };
                });
                setResults(transformedResults);
            })
            .catch((err) => {
                console.log("error: ", err);
            });
    }

    return (
        <React.Fragment>
            <div className="search">
                <Form.Control
                    placeholder="Search by name"
                    onChange={handleChange}
                />
                <ResultsList data={results} />
            </div>
        </React.Fragment>
    );
};

export default Search;
