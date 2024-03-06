// Code generated by smithy-go-codegen DO NOT EDIT.

package cognitoidentity

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Lists all of the Cognito identity pools registered for your account. You must
// use AWS Developer credentials to call this API.
func (c *Client) ListIdentityPools(ctx context.Context, params *ListIdentityPoolsInput, optFns ...func(*Options)) (*ListIdentityPoolsOutput, error) {
	if params == nil {
		params = &ListIdentityPoolsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListIdentityPools", params, optFns, c.addOperationListIdentityPoolsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListIdentityPoolsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// Input to the ListIdentityPools action.
type ListIdentityPoolsInput struct {

	// The maximum number of identities to return.
	//
	// This member is required.
	MaxResults *int32

	// A pagination token.
	NextToken *string

	noSmithyDocumentSerde
}

// The result of a successful ListIdentityPools action.
type ListIdentityPoolsOutput struct {

	// The identity pools returned by the ListIdentityPools action.
	IdentityPools []types.IdentityPoolShortDescription

	// A pagination token.
	NextToken *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListIdentityPoolsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpListIdentityPools{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpListIdentityPools{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "ListIdentityPools"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addOpListIdentityPoolsValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListIdentityPools(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

// ListIdentityPoolsAPIClient is a client that implements the ListIdentityPools
// operation.
type ListIdentityPoolsAPIClient interface {
	ListIdentityPools(context.Context, *ListIdentityPoolsInput, ...func(*Options)) (*ListIdentityPoolsOutput, error)
}

var _ ListIdentityPoolsAPIClient = (*Client)(nil)

// ListIdentityPoolsPaginatorOptions is the paginator options for ListIdentityPools
type ListIdentityPoolsPaginatorOptions struct {
	// The maximum number of identities to return.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListIdentityPoolsPaginator is a paginator for ListIdentityPools
type ListIdentityPoolsPaginator struct {
	options   ListIdentityPoolsPaginatorOptions
	client    ListIdentityPoolsAPIClient
	params    *ListIdentityPoolsInput
	nextToken *string
	firstPage bool
}

// NewListIdentityPoolsPaginator returns a new ListIdentityPoolsPaginator
func NewListIdentityPoolsPaginator(client ListIdentityPoolsAPIClient, params *ListIdentityPoolsInput, optFns ...func(*ListIdentityPoolsPaginatorOptions)) *ListIdentityPoolsPaginator {
	if params == nil {
		params = &ListIdentityPoolsInput{}
	}

	options := ListIdentityPoolsPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListIdentityPoolsPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListIdentityPoolsPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next ListIdentityPools page.
func (p *ListIdentityPoolsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListIdentityPoolsOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxResults = limit

	result, err := p.client.ListIdentityPools(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

func newServiceMetadataMiddleware_opListIdentityPools(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "ListIdentityPools",
	}
}
