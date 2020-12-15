import React from "react";
import { Button, Col, Container, Row } from "react-bootstrap";
import Search from "./Search";

const Home = () => {
    return (
        <div className="home">
            <Container>
                <h1>Project Red String</h1>
                <h3>Some description goes here.</h3>
            </Container>
            <Container>
                <p>Search for a candidate or contributor</p>
                <Search />
                <p>or, browse by category</p>
                <Row>
                    <Col>
                        <Button>City Council Members</Button>
                    </Col>
                    <Col>
                        <Button>Category #2</Button>
                    </Col>
                    <Col>
                        <Button>Category #3</Button>
                    </Col>
                </Row>
            </Container>
        </div>
    );
};

export default Home;
