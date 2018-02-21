import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Button, Form, FormGroup,
  FormControl, ControlLabel, Col, HelpBlock, Checkbox } from 'react-bootstrap';
import { getKeepToken, getTokenGrant } from '../contracts'
import Network from '../network'
import { formatAmount } from '../utils';

const ERRORS = {
  INVALID_AMOUNT: `INVALID_AMOUNT`,
  SERVER: `Sorry, your request cannot be completed at this time.`
};

const RESET_DELAY = 3000; // 3 seconds

class TokenGrantForm extends Component {
  constructor(props) {
    super(props);
    this.state = this.getInitialState();
  }

  getInitialState() {
    return {
      hasError: false,
      requestSent: false,
      requestSuccess: false,
      errorMsg: ERRORS.INVALID_AMOUNT,
    };
  }

  onRequestSuccess() {
    this.setState({
      hasError: false,
      requestSent: true,
      requestSuccess: true
    });
    window.setTimeout(() => {
      this.setState(this.getInitialState());
    }, RESET_DELAY);
  }

  onClick(e) {
    this.submit();
  }

  async submit() {
    const { amount, beneficiary, duration, start, cliff, revocable} = this.state;
    const { tokenGrantContractAddress } = this.props;

    const accounts = await Network.getAccounts();
    const token = await getKeepToken(process.env.REACT_APP_TOKEN_ADDRESS);
    const tokenGrantContract = await getTokenGrant(tokenGrantContractAddress);

    tokenGrantContract.grant(formatAmount(amount, 18), beneficiary, duration, start, cliff, revocable, {from: accounts[0], gas: 90000});
  }

  render() {
    const { amount, beneficiary, duration, start, cliff, revocable,
        hasError,
        errorMsg} = this.state;

    return (
      <div className="token-grant-form">
        <Form horizontal onSubmit={(e) => { e.preventDefault(); }}>
          <FormGroup>
            <Col componentClass={ControlLabel} sm={2}>
              Beneficiary:
            </Col>
            <Col sm={8}>
              <FormControl
                type="text"
                value={beneficiary}
              />
              <HelpBlock className="small">Address to which granted tokens are going to be released.</HelpBlock>
            </Col>
          </FormGroup>

          <FormGroup>
            <Col componentClass={ControlLabel} sm={2}>
              Amount:
            </Col>
            <Col sm={8}>
              <FormControl
                type="text"
                value={amount}
              />
              <HelpBlock className="small">Amount to be granted.</HelpBlock>
            </Col>
          </FormGroup>

          <FormGroup>
            <Col componentClass={ControlLabel} sm={2}>
             Duration:
            </Col>
            <Col sm={8}>
              <FormControl
                type="text"
                value={duration}
              />
              <HelpBlock className="small">Duration in seconds of the period in which the tokens will vest.</HelpBlock>
            </Col>
          </FormGroup>

          <FormGroup>
            <Col componentClass={ControlLabel} sm={2}>
              Start:
            </Col>
            <Col sm={8}>
              <FormControl
                type="text"
                value={start}
              />
              <HelpBlock className="small">Timestamp at which vesting will start.</HelpBlock>
            </Col>
          </FormGroup>

          <FormGroup>
            <Col componentClass={ControlLabel} sm={2}>
              Cliff:
            </Col>
            <Col sm={8}>
              <FormControl
                type="text"
                value={cliff}
              />
              <HelpBlock className="small">Duration in seconds of the cliff after which tokens will begin to vest.</HelpBlock>
            </Col>
          </FormGroup>

          <FormGroup>
            <Col componentClass={ControlLabel} sm={2}>
              Revocable:
            </Col>
            <Col sm={8}>
              <Checkbox validationState="success"></Checkbox>
              <HelpBlock className="small">Whether the token grant is revocable or not.</HelpBlock>
            </Col>
          </FormGroup>

          <Button
            bsStyle="primary"
            bsSize="large"
            onClick={this.onClick.bind(this)}>
            Grant tokens
          </Button>
        </Form>
        { hasError &&
          <small className="error-message">{errorMsg}</small> }
      </div>
    );
  }
}

export default TokenGrantForm;
