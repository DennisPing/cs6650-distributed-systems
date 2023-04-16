# {{classname}}

All URIs are relative to *https://virtserver.swaggerhub.com/IGORTON/Twinder/1.0.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Swipe**](SwipeApi.md#Swipe) | **Post** /swipe/{leftorright}/ | 

# **Swipe**
> Swipe(ctx, body, leftorright)


Swipe left or right

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**SwipeDetails**](SwipeDetails.md)| response details | 
  **leftorright** | **string**| Ilike or dislike user | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

