import React from "react";

class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { error: null, errorInfo: null };
  }

  componentDidCatch(error, errorInfo) {
    // Catch errors in any components below and re-render with error message
    this.setState({
      error,
      errorInfo
    });
  }

  render() {
    const { error, errorInfo } = this.state;
    const { children } = this.props;

    if (!errorInfo) {
      return children;
    }

    return (
      <div>
        <h2>Something went wrong.</h2>
        <pre>
          {error && error.toString()}
          <br />
          {errorInfo.componentStack}
        </pre>
      </div>
    );
  }
}

export default ErrorBoundary;
