package com.example.inventory;

import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

/**
 * REST endpoint for Kubernetes liveness and readiness probes.
 */
@Path("/healthz")
public class HealthResource {

    @GET
    @Produces(MediaType.TEXT_PLAIN)
    public String hello() {
        return "OK";
    }
}
