package com.example.inventory;

import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

@Path("/healthz")
public class HealthResource {

    @GET
    @Produces(MediaType.TEXT_PLAIN)
    public String hello() {
        return "OK";
    }
}
