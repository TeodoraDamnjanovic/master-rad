

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

-- CREATE SEQUENCE public.orders_id_seq
--     START WITH 1
--     INCREMENT BY 1
--     NO MINVALUE
--     NO MAXVALUE
--     CACHE 1;


-- ALTER TABLE public.orders_id_seq OWNER TO postgres2;

-- SET default_tablespace = '';

-- SET default_table_access_method = heap;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres2
--
-- DROP TABLE IF EXISTS public.orders CASCADE;

CREATE TABLE public.orders (
    id integer DEFAULT nextval('public.orders_id_seq'::regclass) NOT NULL,
    user_id integer,
    invoice_id character varying(255),
    created_at timestamp without time zone,
  	paid bool,
  	amount integer,
    status character varying(255)
);


ALTER TABLE public.orders OWNER TO postgres2;

--
-- Name: user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres2
--

SELECT pg_catalog.setval('public.orders_id_seq', 1, true);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres2
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);

-- ALTER TABLE public.orders ADD COLUMN status VARCHAR(255);

INSERT INTO "public"."orders"("user_id","created_at","invoice_id","paid","amount", "status")
VALUES
(100,E'2022-03-14 00:00:00',1,true,1000,'created');
















